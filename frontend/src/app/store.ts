import { combineReducers, configureStore } from "@reduxjs/toolkit";

import themeReducer from "@/app/store/themeSlice";
import userReducer from "@/app/store/userSlice";

const rootReducer = combineReducers({
  theme: themeReducer,
  user: userReducer,
});

export const store = configureStore({
  reducer: rootReducer,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
