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
