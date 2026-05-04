import type { ReactNode } from "react";
import { createContext, useContext, useEffect, useState } from "react";

import ToolBox from "@/components/internal/ToolBox";

export type ToolboxItem = {
  id: string;
  label: string;
  icon: ReactNode;
  content: ReactNode;
};

type ToolboxContextValue = {
  activeItem: ToolboxItem | null;
  activeItemId: string | null;
  closeToolbox: () => void;
  isOpen: boolean;
  items: ToolboxItem[];
  openToolbox: () => void;
  setActiveItem: (itemId: string) => void;
};

const ToolboxContext = createContext<ToolboxContextValue | null>(null);

type ToolboxProviderProps = {
  children: ReactNode;
  items: ToolboxItem[];
};

export function ToolboxProvider({ children, items }: ToolboxProviderProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [activeItemId, setActiveItemId] = useState<string | null>(items[0]?.id ?? null);

  useEffect(() => {
    if (items.length === 0) {
      setActiveItemId(null);
      setIsOpen(false);
      return;
    }

    const activeItemStillExists = items.some((item) => item.id === activeItemId);

    if (!activeItemStillExists) {
      setActiveItemId((items[0] as ToolboxItem).id);
    }
  }, [activeItemId, items]);

  const activeItem = items.find((item) => item.id === activeItemId) ?? null;

  return (
    <ToolboxContext.Provider
      value={{
        activeItem,
        activeItemId,
        closeToolbox: () => setIsOpen(false),
        isOpen,
        items,
        openToolbox: () => {
          if (items.length === 0) {
            return;
          }

          setIsOpen(true);
        },
        setActiveItem: (itemId: string) => setActiveItemId(itemId),
      }}
    >
      {children}
      <ToolBox />
    </ToolboxContext.Provider>
  );
}

export function useToolbox() {
  const context = useContext(ToolboxContext);

  if (!context) {
    throw new Error("useToolbox must be used within a ToolboxProvider.");
  }

  return context;
}
