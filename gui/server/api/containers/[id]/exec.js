import dockerService from '@/server/services/docker';

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const { id } = event.context.params;
    const body = await readBody(event)

    const command = body.command;

    if (!command) {
        return { status: 'error', message: 'Command is required' };
    }

    try {
        const container = await dockerService.getContainer(id);
        const exec = await container.exec({
            AttachStdout: true,
            AttachStderr: true,
            Tty: false,
            Cmd: ['sh', '-c', command]
        });

        const stream = await exec.start({ hijack: true });

        let output = '';
        stream.on('data', (chunk) => {
            output += chunk.toString();
        });

        return new Promise((resolve, reject) => {
            stream.on('end', () => {
                resolve({ output });
            });

            stream.on('error', (error) => {
                reject({ status: 'error', message: error.message });
            });
        });
    } catch (error) {
        return { status: 'error', message: error.message };
    }
});
