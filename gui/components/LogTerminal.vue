<template>
  <div ref="terminal" class="terminal-container"></div>
</template>

<script setup>
import { isClient } from '@vueuse/core';
import { socket } from './socket';

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
    xterm.write(data);
  }
};

const handleShellInput = (data) => {
  socket.emit('shellInput', data);
};

onMounted(async () => {
  if (isClient) {
    const { Terminal } = await import('xterm');
    const { FitAddon } = await import('xterm-addon-fit');
    await import('xterm/css/xterm.css');

    xterm = new Terminal({
      cursorBlink: true,
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
      xterm.write('\r\nConnected to server\r\n');
      connectSocket();
    });
    socket.on('disconnect', () => {
      xterm.write('\r\nDisconnected from server\r\n');
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
  width: 100%;
  padding: 1rem;
  border-radius: 8px;
  background-color: rgb(0, 0, 0);
}

.xterm-viewport {
  height: 73vh;
  overflow-y: auto;

  /* Hide scrollbar */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.xterm-viewport::-webkit-scrollbar {
  display: none;
}
</style>
