<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

final class MerchantReview
{
    private function __construct(
		private MerchantReviewUUID $uuid,
		private string $comment,
		private int $rating,
    ) {}

    public static function create(
        MerchantReviewUUID $uuid,
		string $comment,
		int $rating,
    ): self {
        return new self(
            $uuid,
            $comment,
            $rating,
        );
    }

    public function getUUID(): MerchantReviewUUID
    {
        return $this->uuid;
    }

    public function getComment(): string
    {
        return $this->comment;
    }

    public function getRating(): int
    {
        return $this->rating;
    }
}
