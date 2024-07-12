import { d as defineEventHandler } from '../../../runtime.mjs';
import { d as dockerService } from '../../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'fs';
import 'path';
import 'node:fs';
import 'node:url';
import 'dockerode';

const create = defineEventHandler(async (event) => {
  if (event.request.method !== "POST") {
    return { status: "error" };
  }
  const body = await useBody(event);
  await dockerService.createContainer(body);
  return { status: "created" };
});

export { create as default };
//# sourceMappingURL=create.mjs.map
