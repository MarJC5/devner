import dockerService from '~/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'GET') {
        return { status: 'error', message: 'Invalid request method' };
    }

    const { id } = event.context.params; // Assuming container ID is passed as a route parameter

    try {
        const stats = await dockerService.getContainerStats(id);
        return { status: 'success', data: stats };
    } catch (error) {
        return { status: 'error', message: error.message };
    }
});
