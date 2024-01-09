import { ref } from "vue";
import { ClViewGroup } from "../types";
import { useParent } from "/@/v";

export function useViewGroup() {
	const ViewGroup = ref<ClViewGroup>();
	useParent("cl-view-group", ViewGroup);
	return { ViewGroup };
}
