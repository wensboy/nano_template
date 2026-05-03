import { ReactNode } from "react";

export type LoginFormState = {
  username: string;
  password: string;
  rememberPassword: boolean;
};

export type RegisterFormState = {
  username: string;
  password: string;
  confirmPassword: string;
};

export type InputFieldProps = {
  icon: ReactNode;
  name: string;
  placeholder: string;
  type?: "email" | "password" | "text";
  value: string;
  onChange: (value: string) => void;
};