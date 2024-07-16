import dockerService from '~/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method === 'GET') {
        const containers = await dockerService.getContainers();
        return containers;
    }
});
