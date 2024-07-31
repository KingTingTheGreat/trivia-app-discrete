import dynamic from "next/dynamic";

export default function PlayerPage() {
  const NoSsrPlayerContent = dynamic(
    () => import("@/components/playerContent"),
    {
      ssr: false,
    }
  );
  return <NoSsrPlayerContent />;
}
