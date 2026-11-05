// A ticket's handle: the last six characters of its id, as the CLI prints it.
export const handleOf = (id: string) => id.slice(-6);
