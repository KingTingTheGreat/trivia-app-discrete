import Leaderboard from "@/components/leaderboard";
import BuzzedIn from "@/components/buzzedIn";

const BothPage = () => {
  return (
    <div className="flex justify-around">
      <Leaderboard />
      <BuzzedIn />
    </div>
  );
};

export default BothPage;
