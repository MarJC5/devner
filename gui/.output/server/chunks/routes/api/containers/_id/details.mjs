import { d as defineEventHandler } from '../../../../runtime.mjs';
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

const details = defineEventHandler(async (event) => {
  if (event.node.req.method !== "GET") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  const container = await dockerService.getContainerDetail(id);
  return container;
});

export { details as default };
//# sourceMappingURL=details.mjs.map
