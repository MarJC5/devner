import dockerService from '@/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const body = await useBody(event);
    await dockerService.createContainer(body);
    return { status: 'created' };
});
