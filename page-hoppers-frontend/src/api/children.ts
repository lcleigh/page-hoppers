export interface CreateChildParams {
  name: string;
  age: number;
  pin: string;
  token: string;
}

export async function createChild(params: CreateChildParams) {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const res = await fetch(`${API_URL}/api/children`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${params.token}`,
    },
    body: JSON.stringify({
      name: params.name,
      age: params.age,
      pin: params.pin,
    }),
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.message || 'Failed to add child');
  }
  return res.json();
} 