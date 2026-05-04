import { LuWrench } from "react-icons/lu";

type ToolboxButtonProps = {
  onClick: () => void;
};

export default function ToolboxButton({ onClick }: ToolboxButtonProps) {
  return (
    <button
      aria-label="Open toolbox"
      className="fixed right-6 bottom-6 z-40 flex h-14 w-14 items-center justify-center rounded-full border border-white/20 bg-black/85 text-white shadow-[0_18px_45px_rgba(0,0,0,0.35)] transition duration-200 hover:scale-105 hover:bg-black active:scale-95"
      onClick={onClick}
      type="button"
    >
      <LuWrench size={22} />
    </button>
  );
}

