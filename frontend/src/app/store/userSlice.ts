import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type UserState = {
  user_id: number;
  username: string;
  role_id: number;
  role_level: number;
};

const initialState: UserState = {
  user_id: 0,
  username: "",
  role_id: 0,
  role_level: 0,
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUser(state, action: PayloadAction<UserState>) {
      state.user_id = action.payload.user_id;
      state.username = action.payload.username;
      state.role_id = action.payload.role_id;
      state.role_level = action.payload.role_level;
    },
    clearUser(state) {
      state.user_id = 0;
      state.username = "";
      state.role_id = 0;
      state.role_level = 0;
    },
  },
});

export const { clearUser, setUser } = userSlice.actions;
export default userSlice.reducer;
