const notifications = ref([]);

export const useNotifications = () => {
  const addNotification = (message, type = "info", duration = 5000) => {
    const id = Date.now();
    notifications.value.push({
      id,
      message,
      type, // 'success', 'error',
    });

    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id);
      }, duration);
    }

    return id;
  };

  const removeNotification = (id) => {
    notifications.value = notifications.value.filter((n) => n.id !== id);
  };

  const notifySuccess = (message, duration) => addNotification(message, "success", duration);
  const notifyError = (message, duration) => addNotification(message, "error", duration);
  const notifyInfo = (message, duration) => addNotification(message, "info", duration);
  const notifyWarning = (message, duration) => addNotification(message, "warning", duration);

  return {
    notifications,
    addNotification,
    removeNotification,
    notifySuccess,
    notifyError,
    notifyInfo,
    notifyWarning,
  };
};
