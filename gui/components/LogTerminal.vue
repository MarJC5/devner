<template>
  <div ref="terminal" class="terminal-container"></div>
</template>

<script setup>
import { isClient } from '@vueuse/core';
import { socket } from '~/utils/socket';

const props = defineProps({
  containerId: {
    type: String,
    required: true,
  },
});

const terminal = ref(null);
let xterm;
let fitAddon;

const connectSocket = () => {
  socket.emit('startShell', props.containerId);
};

const disconnectSocket = () => {
  socket.emit('stopShell', props.containerId);
  socket.off('shellOutput', handleShellOutput);
};

const handleShellOutput = (data) => {
  if (xterm) {
    xterm.write(data.toString('utf8'));
  }
};

const handleShellInput = (data) => {
  socket.emit('shellInput', data.toString('utf8'));
};

onMounted(async () => {
  if (isClient) {
    const { Terminal } = await import('xterm');
    const { FitAddon } = await import('xterm-addon-fit');
    await import('xterm/css/xterm.css');

    xterm = new Terminal({
      cursorBlink: true,
      fontSize: 16,
    });
    fitAddon = new FitAddon();
    xterm.loadAddon(fitAddon);

    xterm.open(terminal.value);
    fitAddon.fit();

    xterm.onData(handleShellInput);

    connectSocket();

    socket.on('shellOutput', handleShellOutput);
    socket.on('error', (error) => {
      xterm.write(`\r\nError: ${error}\r\n`);
    });
    socket.on('connect', () => {
      xterm.write('\r\nConnected to ~/server\r\n');
      connectSocket();
    });
    socket.on('disconnect', () => {
      xterm.write('\r\nDisconnected from ~/server\r\n');
    });

    window.addEventListener('resize', fitAddon.fit);
  }
});

onBeforeUnmount(() => {
  if (xterm) {
    disconnectSocket();
    window.removeEventListener('resize', fitAddon.fit);
  }
});
</script>


<style>
.terminal-container {
  height: 75vh;
  max-height: 75vh;
  width: 100%;
  border-radius: 8px;
  background-color: rgb(0, 0, 0);
}

.xterm-screen {
  border-radius: 8px;
  padding: 1rem !important;
}

.xterm-viewport {
  border-radius: 8px;
  overflow-y: auto;

  /* Hide scrollbar */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.xterm-viewport::-webkit-scrollbar {
  display: none;
}
</style>../utils/socket
~/app/utils/socket