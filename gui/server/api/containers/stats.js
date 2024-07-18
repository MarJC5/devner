import dockerService from '~/server/services/docker';

const calculateCPUUsage = (stats) => {
    const cpuDelta = (stats.cpu_stats.cpu_usage.total_usage || 0) - (stats.precpu_stats.cpu_usage.total_usage || 0);
    const systemDelta = (stats.cpu_stats.system_cpu_usage || 0) - (stats.precpu_stats.system_cpu_usage || 0);

    if (systemDelta > 0 && cpuDelta > 0) {
        const cpuUsage = (cpuDelta / systemDelta) * (stats.cpu_stats.online_cpus || 1) * 100;
        return cpuUsage;
    }
    return 0;
};

const aggregateStats = (statsArray) => {
    const aggregated = {
        cpu: 0,
        memory: 0,
        memoryLimit: 0,
        memoryRSS: 0,
        memoryCache: 0,
        memorySwap: 0,
        network: {
            rx: 0,
            tx: 0,
            rxErrors: 0,
            txErrors: 0,
            rxPackets: 0,
            txPackets: 0
        },
        diskIO: {
            readBytes: 0,
            writeBytes: 0,
            readOps: 0,
            writeOps: 0
        },
        pids: 0,
        cpuThrottling: {
            periods: 0,
            throttledPeriods: 0,
            throttledTime: 0
        }
    };

    statsArray.forEach(stats => {
        aggregated.cpu += calculateCPUUsage(stats);
        aggregated.memory += stats.memory_stats?.usage || 0;
        aggregated.memoryLimit += stats.memory_stats?.limit || 0;
        aggregated.memoryRSS += stats.memory_stats?.stats?.rss || 0;
        aggregated.memoryCache += stats.memory_stats?.stats?.cache || 0;
        aggregated.memorySwap += stats.memory_stats?.stats?.mapped_file || 0;
        aggregated.pids += stats.pids_stats?.current || 0;
        aggregated.cpuThrottling.periods += stats.cpu_stats?.throttling_data?.periods || 0;
        aggregated.cpuThrottling.throttledPeriods += stats.cpu_stats?.throttling_data?.throttled_periods || 0;
        aggregated.cpuThrottling.throttledTime += stats.cpu_stats?.throttling_data?.throttled_time || 0;

        if (stats.networks) {
            Object.values(stats.networks).forEach(network => {
                aggregated.network.rx += network.rx_bytes || 0;
                aggregated.network.tx += network.tx_bytes || 0;
                aggregated.network.rxErrors += network.rx_errors || 0;
                aggregated.network.txErrors += network.tx_errors || 0;
                aggregated.network.rxPackets += network.rx_packets || 0;
                aggregated.network.txPackets += network.tx_packets || 0;
            });
        }

        if (stats.blkio_stats && stats.blkio_stats.io_service_bytes_recursive) {
            stats.blkio_stats.io_service_bytes_recursive.forEach(io => {
                if (io.op === 'Read') {
                    aggregated.diskIO.readBytes += io.value || 0;
                }
                if (io.op === 'Write') {
                    aggregated.diskIO.writeBytes += io.value || 0;
                }
            });
        }

        if (stats.blkio_stats && stats.blkio_stats.io_serviced_recursive) {
            stats.blkio_stats.io_serviced_recursive.forEach(io => {
                if (io.op === 'Read') {
                    aggregated.diskIO.readOps += io.value || 0;
                }
                if (io.op === 'Write') {
                    aggregated.diskIO.writeOps += io.value || 0;
                }
            });
        }
    });

    return aggregated;
};

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'GET') {
        return { status: 'error', message: 'Invalid request method' };
    }

    try {
        const stats = await dockerService.getAllContainerStats();
        const aggregatedStats = aggregateStats(stats);

        return { status: 'success', data: aggregatedStats };
    } catch (error) {
        return { status: 'error', message: error.message };
    }
});
