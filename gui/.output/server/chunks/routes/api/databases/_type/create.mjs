import { d as defineEventHandler, r as readBody } from '../../../../runtime.mjs';
import { C as Container, D as Database } from '../../../../_/Database.mjs';
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
  const { type } = event.context.params;
  const { username, password, database } = body;
  const container = await Container.fetchContainerByName(`${type}_devner`);
  if (type === "mysql") {
    await container.cmd(`mysql -u ${username} -p'${password}' -e "CREATE DATABASE IF NOT EXISTS ${database}"`);
  } else if (type === "postgresql") {
    await container.cmd(`psql -U ${username} -c "CREATE DATABASE ${database}"`);
  } else {
    return { status: "error", message: "Invalid database type" };
  }
  const db = await Database.fetchDatabase(database);
  if (!db) {
    return { status: "error", message: "Failed to create database" };
  }
  return { status: "success" };
});

export { create as default };
//# sourceMappingURL=create.mjs.map
