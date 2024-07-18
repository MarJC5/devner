<template>
  <UPageGrid>
    <UPageCard>
      <VisXYContainer>
        <VisLine
          :data="cpuData"
          :x="(d) => new Date(d.time)"
          :y="(d) => d.value"
        />
        <VisAxis
          type="x"
          label="Time"
          :tickFormat="
            (value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)
          "
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
        />
        <VisAxis
          type="x"
          label="Time"
          :tickFormat="
            (value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)
          "
        />
        <VisAxis type="y" label="Memory Usage (MB)" />
      </VisXYContainer>
    </UPageCard>
    <UPageCard>
      <VisXYContainer>
        <VisLine
          :data="networkData"
          :x="(d) => new Date(d.time)"
          :y="(d) => d.value"
        />
        <VisAxis
          type="x"
          label="Time"
          :tickFormat="
            (value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)
          "
        />
        <VisAxis type="y" label="Network Usage (MB)" />
      </VisXYContainer>
    </UPageCard>
    <UPageCard>
      <VisXYContainer>
        <VisLine
          :data="diskIOData"
          :x="(d) => new Date(d.time)"
          :y="(d) => d.value"
        />
        <VisAxis
          type="x"
          label="Time"
          :tickFormat="
            (value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)
          "
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
        />
        <VisAxis
          type="x"
          label="Time"
          :tickFormat="
            (value) =>
              Intl.DateTimeFormat('fr-CH', getTimeFormat(filter)).format(value)
          "
        />
        <VisAxis type="y" label="PIDs" />
      </VisXYContainer>
    </UPageCard>
  </UPageGrid>
</template>

<script setup>
import { VisXYContainer, VisLine, VisAxis, VisSingleContainer, VisDonut, VisBulletLegend } from "@unovis/vue";

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
const networkData = ref([]);
const diskIOData = ref([]);
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
  networkData.value = newStats.map((stat) => ({
    time: stat.timestamp,
    value: (stat.network.rx + stat.network.tx) / (1024 * 1024),
  })); // Convert to MB
  diskIOData.value = newStats.map((stat) => ({
    time: stat.timestamp,
    value: (stat.diskIO.readBytes + stat.diskIO.writeBytes) / (1024 * 1024),
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

<style scoped>
.chart {
  width: 100%;
  height: 400px;
  margin: 20px 0;
}
</style>
