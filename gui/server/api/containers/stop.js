import dockerService from '~/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    await dockerService.stopAllContainers();
    return { status: 'stopped' };
});
