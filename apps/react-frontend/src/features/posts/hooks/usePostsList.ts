import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { getV1Posts } from "@/generated/sdk.gen";
import type { PostListResponse } from "@/generated/types.gen";
import { getToken } from "@/utils/tokenStorage";

interface UsePostsListResult {
  posts: PostListResponse["posts"] | undefined;
  total: number | undefined;
  isLoading: boolean;
  isError: boolean;
}

export function usePostsList(): UsePostsListResult {
  const navigate = useNavigate();

  // Redirect to login if no token is present on mount.
  // getToken() reads localStorage synchronously and is not reactive, so it is
  // intentionally omitted from the dependency array.
  useEffect(() => {
    if (getToken() === null) {
      void navigate({ to: "/login" });
    }
  }, [navigate]);

  const { data, isLoading, isError } = useQuery({
    queryKey: ["posts", "list"],
    enabled: getToken() !== null,
    queryFn: async () => {
      const token = getToken();
      const { data, error } = await getV1Posts({
        baseUrl: import.meta.env.VITE_API_BASE_URL,
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      });
      if (error || !data) throw new Error("Failed to fetch posts");

      return data;
    },
  });

  return {
    posts: data?.posts,
    total: data?.total,
    isLoading,
    isError,
  };
}
