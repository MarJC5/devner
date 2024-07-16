import dockerService from '~/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const { id } = event.context.params;
    const body = await readBody(event)
    const options = body.options;
    if (options) {
        return await dockerService.getContainerLogs(id, options);
    } else {
        return await dockerService.getContainerLogs(id);
    }
});
