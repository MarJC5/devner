import dockerService from '@/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'GET') {
        return { status: 'error' };
    }

    const { id } = event.context.params;
    const container = await dockerService.getContainerDetail(id);
    return container;
});
