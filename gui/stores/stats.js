import { defineStore } from 'pinia';

export const useStatsStore = defineStore('stats', {
    state: () => ({
        stats: localStorage.getItem('dockerStats') ? JSON.parse(localStorage.getItem('dockerStats')) : [],
        lastUpdated: localStorage.getItem('dockerStatsLastUpdated') ? new Date(localStorage.getItem('dockerStatsLastUpdated')) : null,
        historyLength: localStorage.getItem('dockerStatsHistoryLength') || 'month', // default to keeping stats for a month
        filter: localStorage.getItem('dockerStatsFilter') || 'hour', // default to showing all stats
    }),
    getters: {
        lastAggregatedStats(state) {
            if (!state.stats.length) return null;
            return state.stats[state.stats.length - 1];
        },
        currentStats(state) {
            return state.stats[state.stats.length - 1] || {};
        },
    },
    actions: {
        async fetchStats() {
            try {
                const response = await fetch('/api/containers/stats');
                const result = await response.json();
                if (result.status === 'success') {
                    const newStats = result.data;
                    newStats.timestamp = new Date();
                    this.stats.push(newStats);
                    this.lastUpdated = new Date();

                    // Prune old stats based on history length
                    this.pruneStats();

                    // Save to local storage
                    localStorage.setItem('dockerStats', JSON.stringify(this.stats));
                    localStorage.setItem('dockerStatsLastUpdated', this.lastUpdated);
                } else {
                    console.error('Failed to fetch stats:', result.message);
                }
            } catch (error) {
                console.error('Failed to fetch stats:', error);
            }
        },
        pruneStats() {
            const now = new Date();
            let cutoffDate;

            if (!this.stats) {
                return;
            }

            switch (this.historyLength) {
                case 'hour':
                    cutoffDate = new Date(now.setHours(now.getHours() - 1));
                    break;
                case 'day':
                    cutoffDate = new Date(now.setDate(now.getDate() - 1));
                    break;
                case 'week':
                    cutoffDate = new Date(now.setDate(now.getDate() - 7));
                    break;
                case 'month':
                    cutoffDate = new Date(now.setMonth(now.getMonth() - 1));
                    break;
                case 'year':
                    cutoffDate = new Date(now.setFullYear(now.getFullYear() - 1));
                    break;
                case 'all':
                default:
                    return; // Do not prune
            }

            this.stats = this.stats.filter(stat => new Date(stat.timestamp) >= cutoffDate) || [];
        },
        filterStats() {
            const now = new Date();
            let cutoffDate;

            if (!this.stats) {
                return [];
            }

            switch (this.filter) {
                case 'hour':
                    cutoffDate = new Date(now.setHours(now.getHours() - 1));
                    break;
                case 'day':
                    cutoffDate = new Date(now.setDate(now.getDate() - 1));
                    break;
                case 'week':
                    cutoffDate = new Date(now.setDate(now.getDate() - 7));
                    break;
                case 'month':
                    cutoffDate = new Date(now.setMonth(now.getMonth() - 1));
                    break;
                case 'year':
                    cutoffDate = new Date(now.setFullYear(now.getFullYear() - 1));
                    break;
                case 'all':
                default:
                    return this.stats; // Do not filter
            }

            return this.stats.filter(stat => new Date(stat.timestamp) >= cutoffDate) || [];
        },
        setFilter(filter) {
            this.filter = filter || 'hour';
            localStorage.setItem('dockerStatsFilter', filter);
        },
        setHistoryLength(length) {
            this.historyLength = length || 'month';
            localStorage.setItem('dockerStatsHistoryLength', length);
            this.pruneStats();
            localStorage.setItem('dockerStats', this.stats ? JSON.stringify(this.stats) : '[]');
        },
        startTracking(interval = 60) {
            this.fetchStats(); // Fetch immediately
            this.trackingInterval = setInterval(this.fetchStats, interval * 1000);
        },
        stopTracking() {
            clearInterval(this.trackingInterval);
        },
    },
});
