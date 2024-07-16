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

const stop = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  await dockerService.stopContainer(id);
  return { status: "stopped" };
});

export { stop as default };
//# sourceMappingURL=stop.mjs.map
