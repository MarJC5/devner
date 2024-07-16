<template>
    <div ref="logContainer" class="h-[75vh] overflow-y-auto mt-4 p-4 rounded-lg bg-gray-100 dark:bg-gray-800">
      <div v-for="log in logs" :key="log" class="log-entry">{{ log }}</div>
    </div>
  </template>
  
  <script setup>
  import { socket } from '~/utils/socket';
  
  const props = defineProps({
    containerId: {
      type: String,
      required: true
    }
  });
  
  const logs = ref([]);
  const logContainer = ref(null);
  
  const connectSocket = () => {
    socket.emit('getLogs', props.containerId);
  };
  
  const handleLog = (log) => {
    logs.value.push(log);
  };
  
  const handleError = (error) => {
    console.error('Socket error:', error);
  };
  
  const handleConnect = () => {
    console.log('Reconnected to ~/server');
    logs.value = []; // Clear logs on reconnect to avoid duplicate entries
    connectSocket();
  };

  const scrollToBottom = () => {
    if (logContainer.value) {
        logContainer.value.scrollTop = logContainer.value.scrollHeight;
    }
    
};
  
  onMounted(() => {
    connectSocket();
  
    socket.on('log', handleLog);
    socket.on('error', handleError);
    socket.on('connect', handleConnect);
    socket.on('disconnect', () => {
      console.log('Disconnected from ~/server');
    });

    scrollToBottom();
  });
  
  onBeforeUnmount(() => {
    socket.off('log', handleLog);
    socket.off('error', handleError);
    socket.off('connect', handleConnect);
  });
  
  watch(logs, () => {
        scrollToBottom();
    });

  watch(() => props.containerId, () => {
    logs.value = []; // Clear logs when containerId changes
    connectSocket();
  });
  </script>
  
  <style>
  .log-entry {
    white-space: pre-wrap; /* Preserve whitespace and line breaks */
    font-family: monospace; /* Use a monospaced font for logs */
  }
  </style>
  ../utils/socket~/app/utils/socket