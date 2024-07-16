import { d as defineEventHandler, a as dockerService } from '../../../../runtime.mjs';
import 'node:http';
import 'node:https';
import 'events';
import 'https';
import 'http';
import 'net';
import 'tls';
import 'crypto';
import 'stream';
import 'url';
import 'zlib';
import 'buffer';
import 'fs';
import 'path';
import 'engine.io';
import 'socket.io';
import 'dockerode';
import 'node:fs';
import 'node:url';

const restart = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  await dockerService.restartContainer(id);
  return { status: "restart" };
});

export { restart as default };
//# sourceMappingURL=restart.mjs.map
