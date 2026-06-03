import { ReactNode } from "react";

export type LoginFormState = {
  username: string;
  password: string;
  rememberPassword: boolean;
  role?: string;
};

export type RegisterFormState = {
  username: string;
  password: string;
  confirmPassword: string;
  role?: string;
};

export type InputFieldProps = {
  icon: ReactNode;
  name: string;
  placeholder: string;
  type?: "email" | "password" | "text";
  value: string;
  onChange: (value: string) => void;
};
