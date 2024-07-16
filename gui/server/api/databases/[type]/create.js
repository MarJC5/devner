import Container from "~/utils/Container.js";
import Database from "~/utils/Database.js";

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const body = await readBody(event)
    const { type } = event.context.params;
    const { username, password, database } = body;

    const container = await Container.fetchContainerByName(`${type}_devner`);
    if (type === 'mysql') {
        await container.cmd(`mysql -u ${username} -p'${password}' -e "CREATE DATABASE IF NOT EXISTS ${database}"`);
    } else if (type === 'postgresql') {
        await container.cmd(`psql -U ${username} -c "CREATE DATABASE ${database}"`);
    } else {
        return { status: 'error', message: 'Invalid database type' };
    }

    // watch if the project has been created
    const db = await Database.fetchDatabase(database);
    
    if (!db) {
        return { status: 'error', message: 'Failed to create database' };
    }

    return { status: 'success' };
});
