import { useCallback, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import {
  LuBotMessageSquare,
  LuFileText,
  LuFolderOpen,
  LuLayoutDashboard,
  LuListTodo,
  LuSettings2,
  LuUserRound,
} from "react-icons/lu";

import { request } from "@/api/client";
import { UserLogout } from "@/api/user/auth";
import { useAppDispatch } from "@/app/hooks";
import { clearUser } from "@/app/store/userSlice";
import { ToolboxProvider, type ToolboxItem } from "@/app/toolbox";
import SaasLayout from "@/components/custom/saas";
import AiChatToolboxPanel from "@/components/internal/AiChatToolboxPanel";
import {
  ToolboxActionsPanel,
  ToolboxEditorPanel,
  ToolboxFilesPanel,
  ToolboxNotesPanel,
  ToolboxOverviewPanel,
  ToolboxProfilePanel,
} from "@/components/internal/ToolBox";

export default function HomePage() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const handleLogout = useCallback(async () => {
    try {
      await request(() => UserLogout());
    } catch {
      // Keep frontend logout resilient even if the backend cookie is already invalid.
    } finally {
      dispatch(clearUser());
      navigate("/login", { replace: true });
    }
  }, [dispatch, navigate]);

  const toolboxItems = useMemo<ToolboxItem[]>(
    () => [
      {
        content: <ToolboxOverviewPanel />,
        icon: <LuLayoutDashboard size={20} />,
        id: "overview",
        label: "Overview",
      },
      {
        content: <AiChatToolboxPanel />,
        icon: <LuBotMessageSquare size={20} />,
        id: "ai-chat",
        label: "AI Chat",
      },
      {
        content: <ToolboxEditorPanel />,
        icon: <LuFileText size={20} />,
        id: "editor",
        label: "Editor",
      },
      {
        content: <ToolboxFilesPanel />,
        icon: <LuFolderOpen size={20} />,
        id: "files",
        label: "Files",
      },
      {
        content: <ToolboxNotesPanel />,
        icon: <LuListTodo size={20} />,
        id: "notes",
        label: "Notes",
      },
      {
        content: <ToolboxActionsPanel onLogout={handleLogout} />,
        icon: <LuSettings2 size={20} />,
        id: "actions",
        label: "Actions",
      },
      {
        content: <ToolboxProfilePanel />,
        icon: <LuUserRound size={20} />,
        id: "profile",
        label: "Profile",
      },
    ],
    [handleLogout],
  );

  return (
    <ToolboxProvider items={toolboxItems} placement="bottom-right">
      <SaasLayout onLogout={handleLogout} />
    </ToolboxProvider>
  );
}
