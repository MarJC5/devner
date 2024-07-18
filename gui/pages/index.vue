<template>
  <div>
    <UPageHeader
      :headline="lastUpdatedMessage"
      class="sticky top-0 bg-white z-10 dark:bg-gray-800 pb-4"
    >
      <template #title>
        <h1 class="flex">
          Metrics
        </h1>
      </template>
    </UPageHeader>
    <UDashboardToolbar class="pl-0 py-4">
      <template #left>
        <div class="flex flex-wrap items-stretch gap-2">
        <USelectMenu 
          class="capitalize"
          v-model="filter" 
          @change="updateFilter"
          icon="i-heroicons-funnel"
          :options="[
            'Hour',
            'Day',
            'Week',
            'Month',
            'Year',
          ]"
          color="white"
        />
        <div class="flex flex-wrap gap-2">
        <UBadge variant="soft" color="blue" :label="'CPU Usage: ' + latestStats.cpu.toFixed(2) + '%'" />
          <UBadge variant="soft" color="green" :label="'Memory Usage: ' + (latestStats.memory / (1024 * 1024)).toFixed(2) + ' MB'" />
          <UBadge variant="soft" color="orange" :label="'Network RX: ' + (latestStats.network.rx / (1024 * 1024)).toFixed(2) + ' MB'" />
          <UBadge variant="soft" color="purple" :label="'Network TX: ' + (latestStats.network.tx / (1024 * 1024)).toFixed(2) + ' MB'" />
          <UBadge variant="soft" color="red" :label="'Disk Read: ' + (latestStats.diskIO.readBytes / (1024 * 1024)).toFixed(2) + ' MB'" />
          <UBadge variant="soft" color="teal" :label="'Disk Write: ' + (latestStats.diskIO.writeBytes / (1024 * 1024)).toFixed(2) + ' MB'" />
          <UBadge variant="soft" color="gray" :label="'PIDs: ' + latestStats.pids" />
        </div>
      </div>
      </template>
    </UDashboardToolbar>
    <StatsGraph 
      :filteredStats="filteredStats" 
      :filter="filter.toLocaleLowerCase()"
      class="mt-4"
    />
  </div>
</template>

<script setup>
import { useStatsStore } from '~/stores/stats';

const statsStore = useStatsStore();
const filter = ref(statsStore.filter.charAt(0).toUpperCase() + statsStore.filter.slice(1));
const lastUpdated = ref(null);

const updateLastUpdatedMessage = () => {
  if (statsStore.lastUpdated) {
    lastUpdated.value = `Last updated: ${statsStore.lastUpdated.toLocaleString()}`;
  } else {
    lastUpdated.value = '';
  }
};

const updateFilter = () => {
  statsStore.setFilter(filter.value.toLocaleLowerCase());
};

const filteredStats = computed(() => statsStore.filterStats());
const latestStats = computed(() => statsStore.currentStats);

onMounted(() => {
  updateLastUpdatedMessage();
});

watch(() => statsStore.lastUpdated, updateLastUpdatedMessage);
watch(() => filter.value, updateFilter);

const lastUpdatedMessage = computed(() => lastUpdated.value);
</script>

<style scoped>
.chart {
  width: 100%;
  height: 400px;
  margin: 20px 0;
}
</style>
