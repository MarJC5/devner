<template>
    <UPageGrid>
      <UPageCard>
        <VisXYContainer>
          <VisLine
            :data="cpuData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            color="rgb(59 130 246 / 0.7)"
          />
          <VisAxis
            type="x"
            label="Time"
            :tickFormat="(value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)"
          />
          <VisAxis type="y" label="CPU Usage (%)" />
        </VisXYContainer>
      </UPageCard>
      <UPageCard>
        <VisXYContainer>
          <VisLine
            :data="memoryData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            color="rgb(34 197 94 / 0.7)"
          />
          <VisAxis
            type="x"
            label="Time"
            :tickFormat="(value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)"
          />
          <VisAxis type="y" label="Memory Usage (MB)" />
        </VisXYContainer>
      </UPageCard>
      <UPageCard>
        <VisXYContainer>
          <VisLine
            :data="networkRxData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            label="Network RX"
            color="rgb(249 115 22 / 0.7)"
          />
          <VisLine
            :data="networkTxData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            label="Network TX"
            color="rgb(168 85 247 / 0.7)"
          />
          <VisAxis
            type="x"
            label="Time"
            :tickFormat="(value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)"
          />
          <VisAxis type="y" label="Network Usage (MB)" />
        </VisXYContainer>
      </UPageCard>
      <UPageCard>
        <VisXYContainer>
          <VisLine
            :data="diskReadData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            label="Disk Read"
            color="rgb(239 68 68 / 0.7)"
          />
          <VisLine
            :data="diskWriteData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            label="Disk Write"
            color="rgb(20 184 166 / 0.7)"
          />
          <VisAxis
            type="x"
            label="Time"
            :tickFormat="(value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)"
          />
          <VisAxis type="y" label="Disk I/O (MB)" />
        </VisXYContainer>
      </UPageCard>
      <UPageCard>
        <VisXYContainer>
          <VisLine
            :data="pidsData"
            :x="(d) => new Date(d.time)"
            :y="(d) => d.value"
            color="rgb(var(--color-gray-500) / 0.7)"
          />
          <VisAxis
            type="x"
            label="Time"
            :tickFormat="(value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)"
          />
          <VisAxis type="y" label="PIDs" />
        </VisXYContainer>
      </UPageCard>
    </UPageGrid>
  </template>
  
  

  <script setup>
  import { VisXYContainer, VisLine, VisAxis } from "@unovis/vue";
  import { ref, watch } from "vue";
  
  const props = defineProps({
    filteredStats: {
      type: Array,
      required: true,
    },
    filter: {
      type: String,
      required: true,
    },
  });
  
  const cpuData = ref([]);
  const memoryData = ref([]);
  const networkRxData = ref([]);
  const networkTxData = ref([]);
  const diskReadData = ref([]);
  const diskWriteData = ref([]);
  const pidsData = ref([]);
  
  const timeFormats = {
    hour: { hour: "2-digit", minute: "2-digit", second: "2-digit" },
    day: { hour: "2-digit", minute: "2-digit" },
    week: { month: "short", day: "numeric" },
    month: { month: "short", day: "numeric" },
    year: { year: "numeric", month: "short" },
    all: { year: "numeric", month: "short", day: "numeric" },
  };
  
  const getTimeFormat = (filter) => timeFormats[filter] || timeFormats.all;
  
  const updateCharts = (newStats) => {
    cpuData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.cpu,
    }));
    memoryData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.memory / (1024 * 1024),
    })); // Convert to MB
    networkRxData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.network.rx / (1024 * 1024),
    })); // Convert to MB
    networkTxData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.network.tx / (1024 * 1024),
    })); // Convert to MB
    diskReadData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.diskIO.readBytes / (1024 * 1024),
    })); // Convert to MB
    diskWriteData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.diskIO.writeBytes / (1024 * 1024),
    })); // Convert to MB
    pidsData.value = newStats.map((stat) => ({
      time: stat.timestamp,
      value: stat.pids,
    }));
  };
  
  watch(
    () => props.filteredStats,
    (newStats) => {
      updateCharts(newStats);
    },
    { immediate: true }
  );
  </script>
  
  