export namespace Time {
  export const Second = 1000;
  export const Minute = 60 * Second;
  export const Hour = 60 * Minute;
  export const Day = 24 * Hour;

  export const formatDuration = (seconds: number) => {
    const days = Math.floor(seconds / Time.Day);
    seconds %= Time.Day;
    const hours = Math.floor(seconds / Time.Hour);
    seconds %= Time.Hour;
    const minutes = Math.floor(seconds / Time.Minute);
    seconds %= Time.Minute;
    return { days, hours, minutes, seconds };
  };
}

export const formatDistanceToNow = (past: Date): string => {
  const now = new Date();
  const delta = Math.floor((now.getTime() - past.getTime()) / 1000); // in seconds
  if (delta < 60) {
    return `${delta} seconds ago`;
  }
  const minutes = Math.floor(delta / 60);
  if (minutes < 60) {
    return `${minutes} minutes ago`;
  }

  const hours = Math.floor(minutes / 60);
  if (hours < 24) {
    return `${hours} hours ago`;
  }

  const days = Math.floor(hours / 24);
  return `${days} days ago`;
};
