import Container from "@/models/Container";
import Project from "@/models/Project";

export default defineEventHandler(async (event) => {
    if (event.node.req.method !== 'POST') {
        return { status: 'error' };
    }

    const { name } = event.context.params;

    const container = await Container.fetchContainerByName("frankenphp_devner");
    await container.cmd(`rm -rf /var/www/html/${name}`);

    // watch if the project has been deleted
    const project = await Project.fetchProject(name);
    
    if (project) {
        return { status: 'error', message: 'Failed to delete project' };
    }

    return { status: 'success' };
});
