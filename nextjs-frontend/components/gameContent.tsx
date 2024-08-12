import { Player } from "@/types";
import { CSSProperties } from "react";

export const TableRow = ({ children, index }: { children: React.ReactNode; index: number }) => {
  return (
    <tr className="border-2" key={"row: " + index}>
      {children}
    </tr>
  );
};

export const TableData = ({ children, style }: { children: any; style?: CSSProperties }) => {
  return (
    <td className="text-center py-2 px-4 border-2 text-lg font-medium" style={style}>
      {children}
    </td>
  );
};

const GameContent = ({
  title,
  headers,
  content,
  mapFunc,
}: {
  title: string;
  headers: string[];
  content: Player[];
  mapFunc: (row: Player, index: number) => React.ReactNode;
}) => {
  console.log(content.map((player, index) => mapFunc(player, index)));
  return (
    <div className="flex flex-col items-center">
      <h2 className="text-5xl font-semibold p-2 m-1">{title}</h2>
      <table className="border-collapse">
        <thead>
          <tr>
            {headers.map((header, index) => (
              <th className="px-6 py-2 border-solid border-2 text-xl" key={"header: " + index}>
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>{content.map((player, index) => mapFunc(player, index))}</tbody>
      </table>
    </div>
  );
};

export default GameContent;
