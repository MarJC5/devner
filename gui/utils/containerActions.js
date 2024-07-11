export const performContainerAction = async (id, action, loadContainer) => {
    try {
      await $fetch(`/api/containers/${id}/${action}`, { method: 'POST' });
      loadContainer();
    } catch (error) {
      console.error(`Failed to ${action} container ${id}:`, error);
    }
  };
  