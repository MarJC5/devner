import { d as defineEventHandler, a as dockerService } from '../../../runtime.mjs';
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

const create = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const body = await useBody(event);
  await dockerService.createContainer(body);
  return { status: "created" };
});

export { create as default };
//# sourceMappingURL=create.mjs.map
