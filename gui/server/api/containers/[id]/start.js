import dockerService from '@/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const { id } = event.context.params;
    await dockerService.startContainer(id);
    return { status: 'started' };
});