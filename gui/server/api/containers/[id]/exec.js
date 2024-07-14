import Docker from 'dockerode';
import { PassThrough } from 'stream';

const docker = new Docker({ socketPath: '/var/run/docker.sock' });

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
        const container = docker.getContainer(id);
        const exec = await container.exec({
            AttachStdout: true,
            AttachStderr: true,
            Tty: false,
            Cmd: ['sh', '-c', command]
        });

        const stream = await exec.start({ hijack: true });

        // Create separate streams for stdout and stderr
        const stdout = new PassThrough();
        const stderr = new PassThrough();
        container.modem.demuxStream(stream, stdout, stderr);

        let stdoutOutput = '';
        let stderrOutput = '';

        stdout.on('data', (chunk) => {
            stdoutOutput += chunk.toString('utf8');
        });

        stderr.on('data', (chunk) => {
            stderrOutput += chunk.toString('utf8');
        });

        return new Promise((resolve, reject) => {
            stream.on('end', () => {
                resolve({
                    stdout: stdoutOutput,
                    stderr: stderrOutput,
                });
            });

            stream.on('error', (error) => {
                reject({ status: 'error', message: error.message });
            });
        });
    } catch (error) {
        return { status: 'error', message: error.message };
    }
});
