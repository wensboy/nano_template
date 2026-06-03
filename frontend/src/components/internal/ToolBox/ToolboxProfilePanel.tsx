import { useEffect, useMemo, useRef, useState, type ChangeEvent, type ReactNode } from "react";
import { LuSave } from "react-icons/lu";

import { request } from "@/api/client";
import { ListRoles, type Role } from "@/api/role/role";
import { UpdatePassword, UpdateUserProfile } from "@/api/user/user";
import { useAppDispatch, useAppSelector } from "@/app/hooks";
import { setUser } from "@/app/store/userSlice";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";

function FieldLabel({ children }: { children: string }) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return <label className={cx("text-xs font-medium", tone.subtleForeground)}>{children}</label>;
}

type ProfileSectionProps = {
  children: ReactNode;
  onSave: () => void;
  title: string;
};

function ProfileSection({ children, onSave, title }: ProfileSectionProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <section className={cx("rounded-lg border p-4", tone.border, tone.subtle)}>
      <div className="mb-4 flex items-center justify-between gap-3">
        <h3 className={cx("text-sm font-semibold", tone.foreground)}>{title}</h3>
        <button
          aria-label={`Save ${title}`}
          className={cx("flex h-8 w-8 items-center justify-center", getThemeButtonClassName(theme))}
          onClick={onSave}
          type="button"
        >
          <LuSave size={14} />
        </button>
      </div>
      <div className="grid gap-3">{children}</div>
    </section>
  );
}

export default function ToolboxProfilePanel() {
  const dispatch = useAppDispatch();
  const theme = useAppSelector((state) => state.theme);
  const user = useAppSelector((state) => state.user);
  const tone = getThemeTone(theme);
  const avatarInputRef = useRef<HTMLInputElement>(null);
  const [avatarPreviewUrl, setAvatarPreviewUrl] = useState("");
  const [nickname, setNickname] = useState(user.username);
  const [signature, setSignature] = useState("");
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [roleId, setRoleId] = useState(String(user.role_id));
  const [roles, setRoles] = useState<Role[]>([]);
  const userIdentifier = String(user.user_id || user.username || "");

  useEffect(() => {
    setNickname((current) => current || user.username);
    setRoleId(String(user.role_id));
  }, [user.role_id, user.username]);

  useEffect(() => {
    return () => {
      if (avatarPreviewUrl) {
        URL.revokeObjectURL(avatarPreviewUrl);
      }
    };
  }, [avatarPreviewUrl]);

  useEffect(() => {
    let isActive = true;

    void request(() => ListRoles({ page_size: 50 }))
      .then((response) => {
        if (isActive) {
          setRoles(response.data.items);
        }
      })
      .catch(() => {
        if (isActive) {
          setRoles([]);
        }
      });

    return () => {
      isActive = false;
    };
  }, []);

  const roleOptions = useMemo(() => {
    if (roles.length > 0) {
      return roles;
    }

    return user.role_id
      ? [
          {
            id: user.role_id,
            name: `Role ${user.role_id}`,
            description: "",
            level: user.role_level,
            state: 0,
            created_at: "",
            updated_at: "",
          },
        ]
      : [];
  }, [roles, user.role_id, user.role_level]);

  const selectedRole = roleOptions.find((role) => String(role.id) === roleId);
  const roleLevel = selectedRole?.level ?? user.role_level;

  function handleAvatarChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];

    if (!file) {
      return;
    }

    setAvatarPreviewUrl(URL.createObjectURL(file));
  }

  function handleSaveProfile() {
    void request(() =>
      UpdateUserProfile({
        nickname,
        signature,
      }),
    ).catch(() => undefined);
  }

  function handleSavePassword() {
    void request(() =>
      UpdatePassword({
        old_password: oldPassword,
        new_password: newPassword,
      }),
    )
      .then(() => {
        setOldPassword("");
        setNewPassword("");
      })
      .catch(() => undefined);
  }

  function handleSaveAccess() {
    dispatch(
      setUser({
        ...user,
        role_id: Number(roleId) || user.role_id,
        role_level: roleLevel,
      }),
    );
  }

  return (
    <div className="space-y-4">
      <ProfileSection onSave={handleSaveProfile} title="Profile">
        <div className="flex items-start gap-4">
          <button
            aria-label="Choose avatar"
            className={cx(
              "flex aspect-square w-[20%] min-w-24 shrink-0 items-center justify-center overflow-hidden rounded-full border text-xl font-semibold transition-colors",
              tone.border,
              tone.surface,
              tone.foreground,
              tone.focusRing,
              tone.secondaryHover,
            )}
            onClick={() => avatarInputRef.current?.click()}
            type="button"
          >
            {avatarPreviewUrl ? (
              <img alt="User avatar" className="h-full w-full object-cover" src={avatarPreviewUrl} />
            ) : (
              nickname.slice(0, 1).toUpperCase() || "U"
            )}
          </button>
          <input
            accept="image/*"
            className="hidden"
            onChange={handleAvatarChange}
            ref={avatarInputRef}
            type="file"
          />

          <div className="grid min-w-0 flex-1 gap-3">
            <div className="grid gap-1.5">
              <FieldLabel>昵称</FieldLabel>
              <input
                className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
                onChange={(event) => setNickname(event.target.value)}
                value={nickname}
              />
            </div>
            <div className="grid gap-1.5">
              <FieldLabel>用户名 / UID</FieldLabel>
              <input
                className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.muted, tone.mutedForeground)}
                disabled
                value={userIdentifier}
              />
            </div>
            <div className="grid gap-1.5">
              <FieldLabel>签名</FieldLabel>
              <textarea
                className={cx(
                  "min-h-20 w-full resize-none rounded-md border px-3 py-2 text-sm",
                  tone.border,
                  tone.surface,
                  tone.foreground,
                  tone.focusRing,
                )}
                onChange={(event) => setSignature(event.target.value)}
                value={signature}
              />
            </div>
          </div>
        </div>
      </ProfileSection>

      <ProfileSection onSave={handleSavePassword} title="Reset Password">
        <div className="grid gap-1.5">
          <FieldLabel>旧密码</FieldLabel>
          <input
            className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
            onChange={(event) => setOldPassword(event.target.value)}
            type="password"
            value={oldPassword}
          />
        </div>
        <div className="grid gap-1.5">
          <FieldLabel>新密码</FieldLabel>
          <input
            className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
            onChange={(event) => setNewPassword(event.target.value)}
            type="password"
            value={newPassword}
          />
        </div>
      </ProfileSection>

      <ProfileSection onSave={handleSaveAccess} title="Access">
        <div className="grid gap-1.5">
          <FieldLabel>用户角色</FieldLabel>
          <select
            className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
            onChange={(event) => setRoleId(event.target.value)}
            value={roleId}
          >
            {roleOptions.length > 0 ? (
              roleOptions.map((role) => (
                <option key={role.id} value={role.id}>
                  {role.name}
                </option>
              ))
            ) : (
              <option value="0">未设置</option>
            )}
          </select>
        </div>
        <div className="grid gap-1.5">
          <FieldLabel>角色等级</FieldLabel>
          <input
            className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.muted, tone.mutedForeground)}
            disabled
            value={roleLevel}
          />
        </div>
      </ProfileSection>
    </div>
  );
}
