import { create } from "zustand";
import { CastMember, CastScreenshot } from "@/types/api";

interface CastFormState {
  isOpen: boolean;
  editingCast: CastMember | null;
  selectedThumbnail: string | null;
  selectedProfileImage: string | null;
  openForm: (cast?: CastMember, screenshot?: CastScreenshot) => void;
  closeForm: () => void;
}

export const useCastFormStore = create<CastFormState>((set) => ({
  isOpen: false,
  editingCast: null,
  selectedThumbnail: null,
  selectedProfileImage: null,
  openForm: (cast, screenshot) =>
    set({
      isOpen: true,
      editingCast: cast || null,
      selectedThumbnail: screenshot
        ? screenshot.thumbnailName
        : cast
          ? cast.thumbnailName
          : undefined,
      selectedProfileImage: screenshot
        ? screenshot.profileImage
        : cast
          ? cast.profileImage
          : undefined,
    }),
  closeForm: () =>
    set({
      isOpen: false,
      editingCast: null,
      selectedThumbnail: null,
      selectedProfileImage: null,
    }),
}));
