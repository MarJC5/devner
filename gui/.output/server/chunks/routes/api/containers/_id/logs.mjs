import { d as defineEventHandler, r as readBody } from '../../../../runtime.mjs';
import { d as dockerService } from '../../../../_/docker.mjs';
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

const logs = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  const body = await readBody(event);
  const options = body.options;
  if (options) {
    return await dockerService.getContainerLogs(id, options);
  } else {
    return await dockerService.getContainerLogs(id);
  }
});

export { logs as default };
//# sourceMappingURL=logs.mjs.map
