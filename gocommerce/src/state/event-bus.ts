export type AppEvent = {
    type: string;
    payload?: unknown;
};

const listeners = new Set<(event: AppEvent) => void>();

export const appEventBus = {
    publish(event: AppEvent) {
        listeners.forEach(listener => listener(event));
    },
    subscribe(listener: (event: AppEvent) => void): () => void {
        listeners.add(listener);
        return () => listeners.delete(listener);
    },
};
