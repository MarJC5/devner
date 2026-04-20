// DevnerPrompt — a minimal floating popup bound to a global hotkey.
//
// Press the hotkey (Cmd+D by default), type a prompt, hit Return. The
// output of `devner agent --yes <prompt>` streams into the panel.
// Press Esc to hide. The process stays resident so subsequent opens
// are instant.
//
// Build: ./build.sh  (produces DevnerPrompt.app)
//
// Rationale: Go can't register a macOS global hotkey without CGO +
// Carbon/Accessibility bindings, and we want to keep the main binary
// CGO-free for trivial cross-compilation. A small Swift shim owns the
// UI layer and shells out to the Go binary for everything else.

import Cocoa
import Carbon.HIToolbox

// MARK: - Config

// Absolute path to the devner binary. Overridable via DEVNER_BIN env
// var so CI / packaged builds don't have to hardcode /usr/local/bin.
let devnerBinary: String = {
    if let env = ProcessInfo.processInfo.environment["DEVNER_BIN"], !env.isEmpty {
        return env
    }
    // Common install paths in order.
    for candidate in ["/usr/local/bin/devner", "/opt/homebrew/bin/devner"] {
        if FileManager.default.isExecutableFile(atPath: candidate) {
            return candidate
        }
    }
    return "/usr/local/bin/devner"
}()

// MARK: - Global hotkey

final class HotkeyController {
    typealias Handler = () -> Void
    private var hotKeyRef: EventHotKeyRef?
    private var handler: Handler?

    // Register Cmd+D. id 1 is an arbitrary app-scoped tag.
    func register(handler: @escaping Handler) {
        self.handler = handler

        var eventType = EventTypeSpec(eventClass: OSType(kEventClassKeyboard),
                                      eventKind: UInt32(kEventHotKeyPressed))

        InstallEventHandler(GetApplicationEventTarget(), { (_, event, userData) -> OSStatus in
            guard let userData = userData else { return noErr }
            let controller = Unmanaged<HotkeyController>.fromOpaque(userData).takeUnretainedValue()
            controller.handler?()
            return noErr
        }, 1, &eventType, UnsafeMutableRawPointer(Unmanaged.passUnretained(self).toOpaque()), nil)

        let hotKeyID = EventHotKeyID(signature: OSType(0x444E5052 /* 'DNPR' */), id: 1)
        // Cmd+D — change keyCode / modifiers here to rebind.
        let keyCode: UInt32 = UInt32(kVK_ANSI_D)
        let modifiers: UInt32 = UInt32(cmdKey)
        RegisterEventHotKey(keyCode, modifiers, hotKeyID, GetApplicationEventTarget(), 0, &hotKeyRef)
    }

    deinit {
        if let ref = hotKeyRef { UnregisterEventHotKey(ref) }
    }
}

// MARK: - Panel

final class PromptPanel: NSPanel {
    override var canBecomeKey: Bool { true }
    override var canBecomeMain: Bool { true }
}

// MARK: - App

final class AppDelegate: NSObject, NSApplicationDelegate, NSTextFieldDelegate {
    var panel: PromptPanel!
    var input: NSTextField!
    var output: NSTextView!
    var spinner: NSProgressIndicator!
    var hotkey = HotkeyController()
    var runningProcess: Process?

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory) // no dock icon, no menu bar
        buildPanel()
        hotkey.register { [weak self] in self?.toggle() }
    }

    // MARK: UI

    func buildPanel() {
        let width: CGFloat = 640
        let height: CGFloat = 320

        let rect = NSRect(x: 0, y: 0, width: width, height: height)
        panel = PromptPanel(contentRect: rect,
                            styleMask: [.titled, .fullSizeContentView, .nonactivatingPanel],
                            backing: .buffered,
                            defer: false)
        panel.isFloatingPanel = true
        panel.level = .floating
        panel.titleVisibility = .hidden
        panel.titlebarAppearsTransparent = true
        panel.isMovableByWindowBackground = true
        panel.hidesOnDeactivate = false
        panel.hasShadow = true
        panel.backgroundColor = NSColor(calibratedWhite: 0.07, alpha: 0.96)
        panel.isReleasedWhenClosed = false
        panel.collectionBehavior = [.canJoinAllSpaces, .fullScreenAuxiliary]

        let container = NSView(frame: rect)

        // Input field
        input = NSTextField(frame: NSRect(x: 20, y: height - 46, width: width - 40, height: 28))
        input.placeholderString = "Ask the agent — e.g. \"list my laravel projects\""
        input.font = NSFont.systemFont(ofSize: 15)
        input.focusRingType = .none
        input.bezelStyle = .roundedBezel
        input.delegate = self
        input.target = self
        input.action = #selector(submit)
        container.addSubview(input)

        // Spinner
        spinner = NSProgressIndicator(frame: NSRect(x: width - 36, y: height - 46, width: 16, height: 16))
        spinner.style = .spinning
        spinner.isDisplayedWhenStopped = false
        spinner.controlSize = .small
        container.addSubview(spinner)

        // Output
        let scroll = NSScrollView(frame: NSRect(x: 20, y: 20, width: width - 40, height: height - 80))
        scroll.hasVerticalScroller = true
        scroll.borderType = .noBorder
        scroll.drawsBackground = false
        output = NSTextView(frame: scroll.bounds)
        output.isEditable = false
        output.isRichText = false
        output.drawsBackground = false
        output.textColor = NSColor(calibratedWhite: 0.92, alpha: 1.0)
        output.font = NSFont.monospacedSystemFont(ofSize: 12, weight: .regular)
        output.autoresizingMask = [.width]
        output.textContainerInset = NSSize(width: 4, height: 4)
        scroll.documentView = output
        scroll.autoresizingMask = [.width, .height]
        container.addSubview(scroll)

        panel.contentView = container
    }

    func toggle() {
        if panel.isVisible {
            hide()
        } else {
            show()
        }
    }

    func show() {
        // Center on the screen currently containing the cursor.
        if let screen = NSScreen.screens.first(where: { NSMouseInRect(NSEvent.mouseLocation, $0.frame, false) })
            ?? NSScreen.main {
            let frame = screen.visibleFrame
            let size = panel.frame.size
            let x = frame.midX - size.width / 2
            let y = frame.midY - size.height / 2 + 120 // bias up
            panel.setFrameOrigin(NSPoint(x: x, y: y))
        }
        NSApp.activate(ignoringOtherApps: true)
        panel.makeKeyAndOrderFront(nil)
        panel.makeFirstResponder(input)
    }

    func hide() {
        cancelRunning()
        panel.orderOut(nil)
    }

    // MARK: Esc to hide

    func control(_ control: NSControl, textView: NSTextView, doCommandBy commandSelector: Selector) -> Bool {
        if commandSelector == #selector(NSResponder.cancelOperation(_:)) {
            hide()
            return true
        }
        return false
    }

    // MARK: Run devner

    @objc func submit() {
        let prompt = input.stringValue.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !prompt.isEmpty else { return }
        input.stringValue = ""
        appendOutput("▶ \(prompt)\n\n", color: NSColor.systemPurple)
        run(prompt: prompt)
    }

    func run(prompt: String) {
        cancelRunning()
        spinner.startAnimation(nil)

        let p = Process()
        p.executableURL = URL(fileURLWithPath: devnerBinary)
        p.arguments = ["agent", "--yes", "--raw", prompt]
        let stdout = Pipe()
        let stderr = Pipe()
        p.standardOutput = stdout
        p.standardError = stderr

        stdout.fileHandleForReading.readabilityHandler = { [weak self] handle in
            let data = handle.availableData
            if data.isEmpty { return }
            if let s = String(data: data, encoding: .utf8) {
                DispatchQueue.main.async { self?.appendOutput(s) }
            }
        }
        stderr.fileHandleForReading.readabilityHandler = { [weak self] handle in
            let data = handle.availableData
            if data.isEmpty { return }
            if let s = String(data: data, encoding: .utf8) {
                DispatchQueue.main.async { self?.appendOutput(s, color: NSColor.systemRed) }
            }
        }
        p.terminationHandler = { [weak self] _ in
            DispatchQueue.main.async {
                self?.spinner.stopAnimation(nil)
                self?.appendOutput("\n")
            }
        }

        do {
            try p.run()
            runningProcess = p
        } catch {
            appendOutput("\nerror: \(error.localizedDescription)\n", color: NSColor.systemRed)
            spinner.stopAnimation(nil)
        }
    }

    func cancelRunning() {
        runningProcess?.terminate()
        runningProcess = nil
        spinner.stopAnimation(nil)
    }

    func appendOutput(_ text: String, color: NSColor? = nil) {
        let attr = NSMutableAttributedString(string: text)
        attr.addAttribute(.font, value: output.font!, range: NSRange(location: 0, length: attr.length))
        attr.addAttribute(.foregroundColor,
                          value: color ?? output.textColor!,
                          range: NSRange(location: 0, length: attr.length))
        output.textStorage?.append(attr)
        output.scrollToEndOfDocument(nil)
    }
}

// MARK: - Boot

let app = NSApplication.shared
let delegate = AppDelegate()
app.delegate = delegate
app.run()
