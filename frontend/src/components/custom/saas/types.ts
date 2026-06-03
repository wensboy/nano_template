import type { IconType } from "react-icons";

export type SaasNavItem = {
  description: string;
  icon: IconType;
  id: string;
  label: string;
};

export type SaasHeaderAction = {
  disabled?: boolean;
  icon: IconType;
  iconClassName?: string;
  id: string;
  label: string;
  onClick?: () => void;
};

export type SaasNotification = {
  body: string;
  icon?: IconType;
  id: string;
  time: string;
  title: string;
};
