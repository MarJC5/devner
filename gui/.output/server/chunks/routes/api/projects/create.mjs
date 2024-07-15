import { d as defineEventHandler, r as readBody } from '../../../runtime.mjs';
import { C as Container } from '../../../_/Container.mjs';
import { P as Project } from '../../../_/Project.mjs';
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

const create = defineEventHandler(async (event) => {
  if (event.node.req.method !== "POST") {
    return { status: "error" };
  }
  const body = await readBody(event);
  const { name, type } = body;
  const container = await Container.fetchContainerByName("frankenphp_devner");
  switch (type) {
    case "laravel":
      await container.cmd(`composer create-project laravel/laravel /var/www/html/${name}`);
      break;
    case "wordpress":
      await container.cmd(`wp core download --path=/var/www/html/${name} --allow-root`);
      break;
  }
  const project = await Project.fetchProject(name);
  if (!project) {
    return { status: "error", message: "Failed to create project" };
  }
  return { status: "success" };
});

export { create as default };
//# sourceMappingURL=create.mjs.map
