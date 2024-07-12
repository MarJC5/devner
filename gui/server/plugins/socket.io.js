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
