import { d as defineEventHandler } from '../../../../runtime.mjs';
import { C as Container } from '../../../../_/Container.mjs';
import { P as Project } from '../../../../_/Project.mjs';
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
import 'strip-ansi';

const _delete = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const { name } = event.context.params;
  const container = await Container.fetchContainerByName("frankenphp_devner");
  await container.cmd(`rm -rf /var/www/html/${name}`);
  const project = await Project.fetchProject(name);
  if (project) {
    return { status: "error", message: "Failed to delete project" };
  }
  return { status: "success" };
});

export { _delete as default };
//# sourceMappingURL=delete.mjs.map
