import { d as defineEventHandler } from '../../../../runtime.mjs';
import { d as dockerService } from '../../../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'fs';
import 'path';
import 'node:fs';
import 'node:url';
import 'dockerode';

const start = defineEventHandler(async (event) => {
  if (event.request.method !== "POST") {
    return { status: "error" };
  }
  const { id } = event.context.params;
  await dockerService.startContainer(id);
  return { status: "started" };
});

export { start as default };
//# sourceMappingURL=start.mjs.map
