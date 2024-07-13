import { Server as Engine } from 'engine.io';
import { Server } from 'socket.io';
import { defineEventHandler } from 'h3';
import Docker from 'dockerode';

export default defineNitroPlugin((nitroApp) => {
  const engine = new Engine();
  const io = new Server({
    cors: {
      origin: "https://localhost", // Update this with your Nuxt.js app URL
      methods: ["GET", "POST"]
    }
  });
  const docker = new Docker({ socketPath: '/var/run/docker.sock' });

  io.bind(engine);

  io.on('connection', (socket) => {

    socket.on('getLogs', async (containerId) => {
      try {
        const container = docker.getContainer(containerId);
        const logStream = await container.logs({
          follow: true,
          stdout: true,
          stderr: true,
          since: 0,
        });

        logStream.on('data', (chunk) => {
          socket.emit('log', chunk.toString('utf8'));
        });

        socket.on('disconnect', () => {
          logStream.destroy();
        });
      } catch (error) {
        socket.emit('error', 'Error getting logs: ' + error.message);
      }
    });


    socket.on('startShell', async (containerId) => {
      try {
        const container = docker.getContainer(containerId);
        const exec = await container.exec({
          Cmd: ['/bin/sh'] ,
          AttachStdin: true,
          AttachStdout: true,
          AttachStderr: true,
          Tty: true,
        });
  
        const stream = await exec.start({ hijack: true, stdin: true });
  
        stream.on('data', (data) => {
          socket.emit('shellOutput', data.toString('utf8'));
        });

        const shellInputHandler = (data) => {
          stream.write(data);
        };

        socket.on('shellInput', shellInputHandler);

        const cleanup = () => {
          stream.end();
          socket.off('shellInput', shellInputHandler);
        };
        
        socket.on('stopShell', cleanup);
        socket.on('disconnect', cleanup);

      } catch (error) {
        socket.emit('error', 'Error starting shell: ' + error.message);
      }
    });

  });

  nitroApp.router.use('/socket.io/', defineEventHandler({
    handler(event) {
      engine.handleRequest(event.node.req, event.node.res);
      event._handled = true;
    },
    websocket: {
      open(peer) {
        const nodeContext = peer.ctx.node;
        const req = nodeContext.req;

        engine.prepare(req);

        const rawSocket = nodeContext.req.socket;
        const websocket = nodeContext.ws;

        engine.onWebSocket(req, rawSocket, websocket);
      }
    }
  }));
});
