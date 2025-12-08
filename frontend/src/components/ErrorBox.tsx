export function ErrorBox({
  error,
  onClose,
}: {
  error: string;
  onClose: () => void;
}) {
  return (
    <div
      className="bg-card/50 border border-red-400 text-red-400 px-4 py-3 rounded relative mb-4"
      role="alert"
    >
      <strong className="font-bold">Error:</strong>
      <span className="block sm:inline ml-2">{error}</span>
      <span
        className="absolute top-0 bottom-0 right-0 px-4 py-3"
        onClick={onClose}
      >
        <svg
          className="fill-current h-6 w-6 text-red-400"
          role="button"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
        >
          <title>Close</title>
          <path d="M14.348 5.652a1 1 0 00-1.414 0L10 8.586 7.066 5.652a1 1 0 10-1.414 1.414L8.586 10l-2.934 2.934a1 1 0 101.414 1.414L10 11.414l2.934 2.934a1 1 0 001.414-1.414L11.414 10l2.934-2.934a1 1 0 000-1.414z" />
        </svg>
      </span>
    </div>
  );
}
