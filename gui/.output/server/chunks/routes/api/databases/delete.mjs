import { d as defineEventHandler, r as readBody } from '../../../runtime.mjs';
import { C as Container, D as Database } from '../../../_/Database.mjs';
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
  const body = await readBody(event);
  const { type, username, password, database } = body;
  const container = await Container.fetchContainerByName(`${type}_devner`);
  await container.cmd(`mysql -u ${username} -p'${password}' -e "DROP DATABASE IF EXISTS ${database}"`);
  const db = await Database.fetchDatabase(database);
  if (db) {
    return { status: "error", message: "Failed to delete database" };
  }
  return { status: "success" };
});

export { _delete as default };
//# sourceMappingURL=delete.mjs.map
